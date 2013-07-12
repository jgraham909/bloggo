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
    * Split comma-separated tags as individual tags on save, remap to comma-sep on update.
    * Update
    * Delete
* View by tags
* Search
* Aliases
* [hallo.js](http://hallojs.org/ "inline editing")(?)
