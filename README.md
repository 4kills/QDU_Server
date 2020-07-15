# QDU Server

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This server allows you to upload pictures to it, using a no-registration account system. 
The server will return a URL to the client with which they can view the uploaded picture online. The link may be shared. 
The clients can also view all of their all-time uploads and see the likes of each picture individually in a gallery. 

This server is built for [QDU](https://github.com/4kills/QDU).

## Prerequisites

You will need [docker](https://www.docker.com/) and [docker-compose](https://docs.docker.com/compose/) installed on your system 
and you should have basic knowledge of how they work. 

You do NOT need anything else (not even go). 

## Configuring your server

The docker-compose file offers a variety of options and is the main configuration file for your server. 
Especially the 'environment' section under 'qdu-server' should be interesting for you.
You will need to enter your **domain** and the **directory** for saving the pictures. 
Please also adjust the latter under 'volumes'.

## Starting up your QDU Server

To start up your very own QDU Server just download this repository and use the docker-compose. 
This will start up a go-builder container that will compile and build the executable and after that a container actually containing and executing the executable
and a database service. 
You don't need go installed on your system!

## Using TLS

If you want to use the (web-)server with TLS for HTTPS (so the clients' browsers won't freak) you may consider using a reverse proxy such as [NGINX](https://www.nginx.com/), 
which will allow you to fairly easy add TLS certificates and encrypt traffic. 

For obtaining free TLS certificates I can recommend [Certbot](https://certbot.eff.org/). 

NGINX and Certbot work pretty well together and there are hundreds of resources on how to set them up. 
