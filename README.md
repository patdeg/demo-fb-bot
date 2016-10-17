# demo-fb-bot
Demo Facebook Bot with Google App Engine in Go


How To Get Started
==================

more detailed information at https://developers.facebook.com/docs/messenger-platform/quickstart

Step 1: Create a Facebook App and Page
---

* Create a Facebook Page, for example, `My Demo Page` (use the `Create Page` item in the top right dropdown menu on your Facebook home page)

* If needed, set up a Facebook Developer account at https://developers.facebook.com

* Add a New App in your Facebook Developer account

* Select Website as a platform

* Enter your name, contact email, and select the Category `Apps for Pages`

Step 2: Create a Page Token
---

* Go to your Developer Dashboard, i.e. browse to https://developers.facebook.com/apps and click on your new app

* On the left navigation bar, click on `+ Add Product` and select `Messenger`

* Under Token Generation, select the Page you just created. 

* In the popup window, ignore the `Submit for Login Review` message for now, and click `Continue as ...`

* Gather your Page Access Token that we will enter in the code in step 4.

Step 3: Setup Webhook
---

* Under Webhooks, click `Setup Webhooks`

* As Callback URL, enter your application URL followed by `/callback`, for example, `https://facebook-bot-demo.appspot.com/callback`

* Enter your own Verify Token, a keyword/password of your chosing. Keep it as we will enter it in the code in step 4.

* Select the following Subscription Fields: `message_deliveries`, `messaging_optins`, `messages`, and `messaging_postback`

* Do NOT click on `Save and Verify` yet. You need to build and deploy your app first

Step 4: Build and deploy your App Engine
---

* Create a new project on Google Cloud console at https://console.cloud.google.com/home/dashboard by selecting `Create project` in the top left dropdown menu

* Gather your project id, for example, `facebook-bot-demo` in this repository. Initially, Project ID is automatically generated and is different that the Project name that you chose. You can change the project ID by clicking on the "Edit" link. Project ID must be unique across all of Google App Engines, so you will need to try a few times before finding one which is not already taken.

* Git clone this repository locally, i.e. `git clone https://github.com/patdeg/demo-fb-bot.git`

* Edit `demo-fb-bot/app.yaml` and change the application name (currenlty `facebook-bot-demo`) with your project id.

* Edit `demo-fb-bot/const.go` and the constants `PAGE_ACCESS_TOKEN` with your Page Access Token and `VERIFY_TOKEN` with your Verify Token

* Deploy your application with `goapp deploy demo-fb-bot`

* Test the application by browsing to your URL `https://[project-id].appspot.com`

Step 5: Finish setting up your Webhook
---

* Go back to your Facebook Developer tab where you started your Webhook

* Click on `Save and Verify`

* Under Webhooks, subscribe your new page with this webhook by selecting it with the dropdown `Select a Page` and clicking Subscribe

Step 6: Test your bot
---

* Go to your Facebook page, in this example `https://www.facebook.com/deglonbotdemo`

* Click on the Message dropdown bellow the Cover

* Enter a message

* The bot should repeat your message with the word "Hello" appended

* Until your submit your bot to Facebook for approval, only selected users can use the app. To add a user, click Roles in the left navigation bar in your Facebook Developer dashboard, then add them as Testers

Step 7: Develop your bot
---

* Now, you can develop the artificial intelligence of your bot in the function `GetResponse` in `facebook.go` which return the response the bot should give for an incomming `message` from user `facebookUser`