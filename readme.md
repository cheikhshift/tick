# Ticker - Live data streams on Golang
![](https://cdn.xl.thumbs.canstockphoto.com/canstock22242451.jpg)
![enter image description here](https://camo.githubusercontent.com/ed1230a48b6946283ee0d76185a726b49ba58254/68747470733a2f2f7472617669732d63692e6f72672f746f6f6c732f676f6465702e737667)

This Gopher Sauce framework requires [Gopher sauce](https://github.com/cheikhshift/Gopher-Sauce) to work as well as a notion of its template engine.

Tick tracks realtime data from your Mongo database to your client's browser. This is possible with the help of web sockets `ws://`


# How to use it
Ticker supports templates by default. To use a Template function  within your `gos` app functions add the prefix `net_<function name>` to your function name.

### Import Ticker framework

	<import src="github.com/cheikhshift/tick/gos.gxml"/>

### Init and set MongoDB
Add this snippet within your `<main>` tag. Copy the whole snippet if you're Gopher Sauce project has no main tag.
				
				//Tick encryption key	
		tick.Key = "example key 1234"
		//HostName string syntax [host]:[port]
		// var HostName string included on tick/gos.gxml import
		HostName = "localhost:8080"
		// dbs included on tick/gos.gxml import
		dbs,_ = db.Connect("localhost", "database")
		tick.SetDb(dbs)

### Declare DB Object pipeline
Struct of `dbObjcet`. This Struct will be made available on `tick` import :

		<struct name="dbObject">
		 	Id bson.ObjectId `bson:"_id,omitempty"`
		 	Created time.Time
		 	Track int
		 	Signature string `valid:"unique"`
		</struct>

Add a new `method` within your `<methods>` tag (in Gopher Sauce). Replace dbObject with your own model. 
This snippet will find any database object :

		<method name="GetdbObject" return="dbObject">
			dbo := dbObject{}
			dbs.Q(&dbo).Find(nil).One(&dbo)
			return dbo
		</method>

### Template

With the prior section in mind, this is how you track your database object's values. Within your templates on Gopher sauce use the following code.

Get object via pipeline : 

	{{ $item := GetdbObject }}

Open socket to item :

	{{ Tick $item }}

(optional) Display live value of object key. (Any struct field must be in lower case) : 

	{{ Watch $item "track" }}

Syntax :
		
	{{ Watch mongo_databse_object key_string }}		
		
	
# Javascript API
Prior to calling each of the JS functions, remember to invoke {{Tick database_object }} to load the API.
Follow this syntax to create event handlers for your object : 


		<script>
     	//Tick Js function
     	Tick({{ IdOf $item }} , function(data){
     			console.log("Update")
     	})
     	</script>

JS function : Tick(id string , callback function(data) )
- id : ID of database item.
- callback : function called on each update.

### Detect websocket close on client :

	<script>
		$tick.close = function(event ) {
			console.log(event)
		}
	</script>