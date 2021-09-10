FROM amazon/aws-lambda-provided:al2

RUN yum update \
&& yum install -y \
  yum-utils \
  wget \
  go

# For more info, visit: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SSL-on-amazon-linux-2.html#letsencrypt

# Download EPEL
RUN wget -r --no-parent -A 'epel-release-*.rpm' https://dl.fedoraproject.org/pub/epel/7/x86_64/Packages/e/ \
&& rpm -Uvh dl.fedoraproject.org/pub/epel/7/x86_64/Packages/e/epel-release-*.rpm \
&& yum-config-manager --enable epel*

RUN mkdir /etc/letsencrypt
COPY ./acme_validation.sh /etc/letsencrypt

# Install Certbot and its dependencies
RUN yum install -y certbot python2-certbot-apache
