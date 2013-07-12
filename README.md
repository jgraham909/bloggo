# bloggo

Revel example blogging application.

This is meant to be an example blogging application for the awesome [revel](github.com/robfig/revel "revel") framework
See the related [discussion](https://groups.google.com/forum/#!topic/revel-framework/-Uy98Bsm4y8 "google groups discussion") in the revel framework users group.

# Configuration

## Database

The database is configurable via the app.conf directive `bloggo.db` if you do not set this it will default to "bloggo"

    bloggo.db = bloggo



## TODO/Roadmap 

*   Article CRUD
    * Split comma-separated tags as individual tags on save, remap to comma-sep on update.
    * Update
    * Delete
* View by tags
* Search
* Aliases
* [hallo.js](http://hallojs.org/ "inline editing")(?)
