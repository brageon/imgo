# Image upload with DNS
Go upload server in EC2 for domain rjtve.com in Cloudflare. This setup use HAProxy for SSL termination where the certificate was assigned by Certbot and uploaded to ACM. In the Console navigate to Application Load Balancer and create 443 Listener. The instance is Ubuntu 22.04 with t3.micro (free-tier). EC2 IP cost $0.005 per hour.
```
certbot certonly --dns-cloudflare --dns-cloudflare-credentials cloudflare.ini -d rjtve.com
openssl x509 -in fullchain.pem -out cert.pem -first-cert
aws acm import-certificate --certificate fileb://cert.pem --private-key fileb://privkey.pem
echo $(cat fullchain.pem) | xclip -selection clipboard
aws elbv2 create-target-group --name "apollo" --target-type instance --vpc-id vpc-08af9d67b0d --port 80 --protocol TCP
aws elbv2 register-targets --target-group arn:aws:elasticloadbalancing:region:root:targetgroup/apollo/b4089270f058 --targets Id=i-instance
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
Use systemctl to define a backend server that is always open. This way users can upload files to S3 without runtime failure from missing files like it would happen in Python. PowerDNS can be used with SQL/REST API. CPanel in VPS is a GUI adapter of this. Data races or birthday attacks can be setup with Siege or Locust.
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
Backup: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/creating-an-ami-ebs.html
