# nuntius
<b>nuntius</b> is a developer command line tool for sending Firebase push notifications to Android devices,
iOS devices, and browsers (web push notifications).  "nuntius" is the Latin word for messenger.  I don't explain here how to 
enable your devices to receive push notificaitons.  Google has tutorials on that here: [Android](https://firebase.google.com/docs/cloud-messaging/android/client),
[iOS](https://firebase.google.com/docs/cloud-messaging/ios/client), and [Web](https://firebase.google.com/docs/cloud-messaging/js/client).
However, once you <em>do</em> have your device enabled to receive notifications, sending notifications is fairly easy with nuntius.
First get your credentials file from Firebase.  <b>./nuntius -help</b> explains how:
```
  -credentialsFile string
    	A Firebase credentials file downloaded from the Firebase console.
    	Log into the Firebase console and go to your project. 
    	Beside "Project Overview" on the left click the gear/settings icon.  
    	Select "Project Settings".  In the Project Settings page, click on "Service accounts".  
    	Scroll down and click on "Generate new private key".  This is your credentials file.
```
After you have your credentials file, create a message file in this format (again explained by <b>./nuntius -help</b>):
```
  -pushFile string
    	A file that contains a notification (title, body, image (optional)) 
    	and/or data (a map of key/value pairs) that will be pushed to the client(s). 
    	Format:
    	{
    	   "notification" : { 
    	      "title" : "Test Title",
    	      "body"  : "Test Body",
    	      "image" : "https://whatever.com/image.png"
    	   },  
    	
    	   "data" : { 
    	     "key1" : "value1",
    	     "key2" : "value2"
    	   }   
    	}
```

Finally, you'll need either the Firebase token for your device or the topic channel that your device is subscribed to.  
The easiest way to get your device token is to log it when the device registers with Firebase. You then send a push notifications like this (token)
:
```
./nuntius -credentialsFile=firebaseProjectCreds.json -pushFile=someMessageFile.json -token=fYYY08QjXXXXX91bFXm9hKq3VjXXXXXBGqAWZQfa6aZYYYYYSoX8-Qrho7nqI1KTtRdOMXXXXXXXXXXXIo0MqwLP8b2FjugpTe07bT3wm5DgXXxUux-
```
or this (topic):
```
./nuntius -credentialsFile=fireaseProjectCreds.json -pushFile=someMessageFile.json -topic=myClientsSubscribedTopic
```

### Command line options
Here's the complete list of command line options.  
```
  -help 
      Print out all the options
  -credentialsFile string
    	A Firebase credentials file downloaded from the Firebase console.
    	Log into the Firebase console and go to your project. 
    	Beside "Project Overview" on the left click the gear/settings icon.  
    	Select "Project Settings".  In the Project Settings page, click on "Service accounts".  
    	Scroll down and click on "Generate new private key".  This is your credentials file.
  -pushFile string
    	A file that contains a notification (title, body, image (optional)) 
    	and/or data (a map of key/value pairs) that will be pushed to the client(s). 
    	Format:
    	{
    	   "notification" : { 
    	      "title" : "Test Title",
    	      "body"  : "Test Body",
    	      "image" : "https://whatever.com/image.png"
    	   },  
    	
    	   "data" : { 
    	     "key1" : "value1",
    	     "key2" : "value2"
    	   }   
    	}
  -token string
    	token to send the message to
  -topic string
    	topic to send the message to
  -ttl int
    	Time-to-live value for notifications in seconds.
    	0 (default) means "now or never", that is,
    	deliver the message now or don't deliver it at all.
    	Max value is 2419200 (28 days).
    	For details see: https://firebase.google.com/docs/cloud-messaging/concept-options
```
