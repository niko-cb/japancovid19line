# japancovid19line
Essentially the backend for a line-bot that allows access to prefecture data for covid19

# LINE
It uses a webhook (LINE → Dialogflow) so if someone sends a message to the line bot, it will automatically check with Dialogflow for the response. 

# Dialogflow SDK
This code utilizes the Dialogflow SDK in order to automate intent creation on Dialogflow
・https://dialogflow.com/

# Datastore
It uses GCP's Datastore for holding the prefecture data after scraping the JSON data
・https://cloud.google.com/datastore

# Cron Jobs
Scraping and Intent creation jobs are automated and initiated with the help of a free cron job service.
・https://cron-job.org

# Heroku service
The application is uploaded to and hosted by Heroku (free version) so that the cron jobs have something to point to.
・https://heroku.com

# Note
I may eventually host it on a free service with GCP. Either GAE or GKE for practice.