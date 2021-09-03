FROM amazon/aws-lambda-provided:al2

RUN yum update \
&& yum install -y \
  yum-utils \
  wget \
  unzip \
  go

# For more info, visit: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SSL-on-amazon-linux-2.html#letsencrypt

WORKDIR /app

# Install aws cli v2
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
&& unzip awscliv2.zip \
&& ./aws/install -i /usr/local/aws-cli -b /usr/local/bin \
&& rm awscliv2.zip \
&& rm -r aws

# Download EPEL
RUN wget -r --no-parent -A 'epel-release-*.rpm' https://dl.fedoraproject.org/pub/epel/7/x86_64/Packages/e/ \
&& rpm -Uvh dl.fedoraproject.org/pub/epel/7/x86_64/Packages/e/epel-release-*.rpm \
&& yum-config-manager --enable epel*

# Install Certbot and its dependencies
RUN yum install -y certbot python2-certbot-apache
