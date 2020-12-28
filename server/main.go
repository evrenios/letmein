package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/bluele/slack"

	"github.com/evrenios/letmein/misc"
	"github.com/labstack/echo"
)

func main() {

	sess, ee := session.NewSession()
	if ee != nil {
		panic(ee)
	}
	config := &aws.Config{Region: aws.String(conf.AWSRegion)}
	ec2Svc := ec2.New(sess, config)

	e := echo.New()

	e.POST("/api/auth/v1/grant", func(c echo.Context) error {

		ts := time.Now().Unix()
		req := &misc.AuthReq{
			Timestamp: ts,
		}

		if err := c.Bind(req); err != nil {
			conf.errChan <- err
			return err
		}
		if req.Secret != misc.Secret {
			err := errors.New("go away")
			conf.errChan <- err
			return err
		}

		if _, ok := conf.restrictedIPs[req.IP]; ok {
			conf.slacker.notifyWithAttachment(
				fmt.Sprintf("UNAUTHORIZED"),
				&slack.Attachment{
					Text:  fmt.Sprintf("USER: `%s` tried to auth restricted IP: `%s`", req.Name, req.IP),
					Color: "#171617",
				})

			return nil
		}

		req.IP = fmt.Sprintf("%s/32", string(req.IP))
		conf.currentSet[ts] = req
		var err error

		for sg, port := range conf.SGPorts {
			_, err = ec2Svc.AuthorizeSecurityGroupIngress(
				&ec2.AuthorizeSecurityGroupIngressInput{
					CidrIp:     aws.String(req.IP),
					GroupId:    aws.String(sg),
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int64(port),
					ToPort:     aws.Int64(port),
				},
			)
			if err != nil {
				conf.errChan <- err
			}
		}
		conf.slacker.notifyWithAttachment(
			fmt.Sprintf("New IP has been Authorized"),
			&slack.Attachment{
				Text:  fmt.Sprintf("User: `%s`\nIP: `%s`\nFor: `%d h`", req.Name, req.IP, req.Hour),
				Color: "#32a852",
			})

		time.AfterFunc(time.Duration(req.Hour)*time.Hour, func() {
			delete(conf.currentSet, ts)
			deleteAccess(ec2Svc, req)
		})
		return nil
	})

	go e.Start(":80")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	// gracefully cleanup temporarly added ip blocks if program stops
	for _, req := range conf.currentSet {
		deleteAccess(ec2Svc, req)
	}

}

func deleteAccess(ec2Svc *ec2.EC2, req *misc.AuthReq) {
	for sg, port := range conf.SGPorts {
		_, err := ec2Svc.RevokeSecurityGroupIngress(
			&ec2.RevokeSecurityGroupIngressInput{
				CidrIp:     aws.String(req.IP),
				GroupId:    aws.String(sg),
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(port),
				ToPort:     aws.Int64(port),
			},
		)
		if err != nil {
			conf.errChan <- err
		}
	}
	conf.slacker.notifyWithAttachment(
		fmt.Sprintf("IP has been Revoked"),
		&slack.Attachment{
			Text:  fmt.Sprintf("User: `%s`\nIP: `%s`", req.Name, req.IP),
			Color: "#c4352d",
		})

}
