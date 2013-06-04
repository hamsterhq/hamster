**Hamster** is a simplistic parse.com like backend as a service(BaaS).It provides a REST interface to datastores for CRUD operations plus some services.

Documentation is coming soon...yep...soon.Though you can look at the test code which pretty much covers all the features implemented.

>*"Writing test code is better than writing documentation"* - Lazy and Procastinating Anonymous Guy.

See: hamster_test.go & `go test`

**Dependencies**

1. Mongodb: default port.
2. Redis: default port.

Check hamster.toml to change mongodb host string with your own username:password.


**Why Mongodb?**

Hamster is primarily focussed on mobile apps which need quick and dirty datastores. Another reason was the awesome Go mongodb driver:
http://labix.org/mgo. Thank you Gustavo Niemeyer!

That said, there is no reason why we can't have sql databases as an option. Open to implementing MySQL, PostgreSQL integration
in the future.


**What is Redis for?**

Caches, message(push notifications) and task(send emails, long running) queues.


**Why reinvent the wheel?**


As far as I know there are not many modern open source BAAS solutions out there. And though my opinion is highly subjective, I think
Go provides an oppurtunity to build something simple, stable, fast and scalable. In any case, it always good to have a free and open source
project at our hands.

**What is the future of this project?**


I am building my start-up's backend around Hamster, so it's going to be around for a while. 

**Features Implemented:**

1. Developer Accounts.
2. Mutlple Apps( with individual api-key and api-secret).
3. CRUD schemaless objects for the apps.
4. Save and Get Files( stored in GridFS).

Please go through `hamster_test.go`. I will try to get to documentation soon.

**Roadmap**

1. Better Authentication
2. Android, iOS clients
3. Push Notifications(ios, android, browser)
4. API metering and throttling
5. Email service(email verification, forgot password etc.)
6. Javascript, Python, Ruby, Java, PHP clients.
7. Analytics.


**Contributions**


Not accepting pull requests right now due to the highly unstable code, but code reviews, feature requests and issues are most welcome. The pull
request freeze is expected to be around for another four weeks. Sorry about that!

**Why is Licence Missing?**

Hamster is going to be free and open source. Still figuring out the right licencse. I want to keep it open source 
and free but maybe sell a commercial support or provide cloud hosting services. Advice on this is most welcome.
