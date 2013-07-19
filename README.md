# bloggo

Revel example blogging application. This requires a mongo db server, the defaults are detailed below. See the documentation for [revmgo](https://github.com/jgraham909/revmgo "revmgo") for more additional configuration options.

This is meant to be an example blogging application for the awesome [revel](https://github.com/robfig/revel "revel") framework
See the related [discussion](https://groups.google.com/forum/#!topic/revel-framework/-Uy98Bsm4y8 "google groups discussion") in the revel framework users group.

# Example users and content

Bloggo ships with two default user accounts and a few articles. You can log in as either 'Jane Doe' or 'John Doe' by using `jane@example.com` or `john@example.com` respectively with the password `12341234` for either user account.

# Configuration

## Database

The database is configurable via the app.conf directive `bloggo.db` if you do not set this it will default to "bloggo"

    bloggo.db = bloggo

## Collections

Revel models correspond to mongo collections. By default the collection name defaults to the model name eg. Article uses the "Article" collection in mongodb. If you want to store the model data in a different collection use the following config one per model type.

    bloggo.db.collection.MODEL_NAME = "COLLECTION_NAME"

eg.

    bloggo.db.collection.Articles = "articles"
    
## Admin User

The admin user can be set via app.conf, with the directive `bloggo.admin` default this ships with it set as 'Jane Doe' one of the default users. The value is set to the mongo _id value of the user account that should be considered as the admin account.

eg. 
    
    bloggo.admin = "51e9aa4049a1b716bb000003"

## TODO/Roadmap

* User Landing pages
* User Management
    * Creating accounts: limit to just the admin account
    * Reset Password
    * Locale
    * Timezone
    * Profile pictures
    * Deleting accounts: limit to just the admin account, and block deleting own account
* Aritcles
    * Published/unpublished
    * Content pager
    * Article Pictures
    * Fix article sorting by date
* Code Review & Refactoring
    * Rework any flakey abstractions
    * Determine how code can be better re-used.
        * Template helper functions
        * App config for date formatting
* Internationalization show examples of UI/UX changing based on user settings
* Search
