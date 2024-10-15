<h1 align="center">Image Upload DNS</h1>
Go upload server in EC2. This setup use HAProxy for SSL termination where the certificate was assigned by Certbot and uploaded to ACM. In the Console navigate to Application Load Balancer and create 443 Listener. The instance is t3.micro on the free tier. EC2 IP cost $0.005 per hour. VPC resources can't be monitored in the "Billing" section. 

```
Download https://github.com/moovweb/gvm
gvm install go1.23 && ps aux | grep gvm
ps -p <pid> -o %cpu,%mem,cmd && top -p
```

```
certbot certonly --dns-cloudflare --dns-cloudflare-credentials cloudflare.ini -d rjtve.com
openssl x509 -in fullchain.pem -out cert.pem -first-cert
aws acm import-certificate --certificate fileb://cert.pem --private-key fileb://privkey.pem
echo $(cat fullchain.pem) | xclip -selection clipboard
aws elbv2 create-target-group --name "apollo" --target-type instance \n
--vpc-id vpc-<uid> --port 80 --protocol TCP
aws elbv2 register-targets --target-group \n
arn:aws:elasticloadbalancing:region:root:targetgroup/apollo/<id> --targets Id=i-instance
aws elbv2 describe-target-health --target-group 
```

Use EC2 Connect Endpoint from VPC instead of SSH. Setup awscli credentials and use S3 to copy and download your files in a private bucket.

```
cat cert.pem privkey.pem > /etc/ssl/certs/cert.pem
systemctl restart nginx haproxy
dig @ns...com rjtve.com, curl -v, curl -vv
ping, nslookup https://www.rjtve.com/
ifconfig
ss -tuln | grep 127.0.0.53
lsof -i:port, pidof nginx, ps -p
netstat -plant | grep 80
go build -o app.go upload.go /var/www/rjtve/app.go
ps aux | grep go
```

Use systemctl to define a backend server that is always open. This way users can upload files to S3 without runtime failure from missing files like it would happen in Python. PowerDNS can be used with SQL/REST API. CPanel in VPS is a GUI adapter of this. Data race and birthday attack can be setup with Siege and/or Locust.

```
nano /etc/systemd/system/imgo.service
[Unit]
Description=My Go App
After=network.target

[Service]
WorkingDirectory=/var/www/rjtve
User=ubuntu  # whoami
ExecStart=/usr/bin/go run app.go :8081
Restart=always

[Install]
WantedBy=multi-user.target

sudo systemctl enable imgo.service
sudo systemctl start imgo.service
(Optional: Enable DNSSEC)
sudo dnssec-keygen -a RSASHA256 -f KSK rjtve.com
sudo dnssec-keygen -a RSASHA256 -b 2048 -n ZONE rjtve.com
Monitor with mxtoolbox, intodns.com, and whatsmydns.net
```

Backup Links: 

[1] https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-store-restore.html#store-ami

[2] https://developer.hashicorp.com/packer/docs/intro/use-cases
