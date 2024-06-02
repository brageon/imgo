server {
    listen 80;  # Listen on port 80 (forwarded from HAProxy)
    server_name rjtve.com www.rjtve.com;

    root /var/www/rjtve;

    index index.html index.htm index.nginx-debian.html;

    location / {
        proxy_pass http://127.0.0.1:53;
        proxy_cache my_cache;
        proxy_cache_valid 200 302 10m;
        proxy_cache_valid 404 1m;
        proxy_cache_methods GET HEAD;
        try_files $uri $uri/ =404;
        expires 1d;
        add_header Cache-Control "public";
    }
}
