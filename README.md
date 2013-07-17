# bloggo

Revel example blogging application.

This is meant to be an example blogging application for the awesome [revel](github.com/robfig/revel "revel") framework
See the related [discussion](https://groups.google.com/forum/#!topic/revel-framework/-Uy98Bsm4y8 "google groups discussion") in the revel framework users group.

# Configuration

## Database

The database is configurable via the app.conf directive `bloggo.db` if you do not set this it will default to "bloggo"

    bloggo.db = bloggo

## Collections

Revel models correspond to mongo collections. By default the collection name defaults to the model name eg. Article uses the "Article" collection in mongodb. If you want to store the model data in a different collection use the following config one per model type.

    bloggo.db.collection.MODEL_NAME = "COLLECTION_NAME"
    
eg.

    bloggo.db.collection.Articles = "articles"

## TODO/Roadmap 

*   Article CRUD
    * Update: Needs some UI/UX cleanup for when the link is (not) displayed, currently it is always displayed.
    * Delete
* UI cleanup
* User Management
    * Manage creation of accounts
    * Reset Password
    * Locale
    * Timezone
* Code Review & Refactoring
    * Rework any flakey abstractions
    * Determine how code can be better re-used.
        * Template helper functions
        * App config for date formatting
* Internationalization show examples of UI/UX changing based on user settings
* Search
