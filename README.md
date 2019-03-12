# Prerequisites
  * You will need a mySQL database meeting following conditions: 
    - A table must be called "pics"
    - The columns/attributes of "pics" must be: 
      + "pic_id" Type Binary(16) Primary Key
      + "token" Type Binary(16) 
      + "timestamp" Type Timestamp default Current_Timestamp
      + "clicks" Type Int Unsigned default 0
    - This database must run on the same system (VM) as your qdu server (improvements comming soon)
  * QDU is depending on several open source go repositories
  * Windows or Linux as OS
    
# Setting up your QDU Server
  * As stated above you will need a mysql database. There is an sql file in the example directory that creates the database 
  with the required table for you.
  * Go will automatically add the dependencies upon go build-ing the project, so you won't have to worry about that.
  * There are several example settings in the example folder. You can have a look at those if you have trouble 
  understanding some of the settings. 
  * You can use any kind of (reverse)proxy with this project's web server even when using TLS/SSL
