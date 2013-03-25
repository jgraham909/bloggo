# bloggo

Revel example blogging application.

This is meant to be an example blogging application for the awesome [revel](github.com/robfig/revel "revel") framework
See the related [discussion](https://groups.google.com/forum/#!topic/revel-framework/-Uy98Bsm4y8 "google groups discussion") in the revel framework users group.

Left to do;

*   Refactor to use models properly for loading, querying & saving. Pass in either the session or the controller
*   User editing
    *  attempt to avoid duplication of code surrounding Create, Update &amp; change password workflows
*   Article CRUD
    * Fields {Title, Tags, Body}
    * Search
    * Comments(?)
    * Aliases(?)
    * [hallo.js](http://hallojs.org/ "inline editing")
