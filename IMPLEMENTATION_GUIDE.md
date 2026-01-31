# Implementation Guide

## General Aspects

 - Quest YAML files are stored in a configurable folder on the server.
 - Sub-folders are used to categorize and manage the quest files.
 - The editor should remember the positions or quest nodes on the
   canvas, so that when the quest is edited again the user doesn't
   need to rearrange the nodes.

## Implementation Requirements

 - Keep external dependencies to a minimum.
 - Don't add huge external dependencies for small amounts of functionality. Instead,
   implemented the functionality within this project.
 - Code should be readable and understandable by human developers.
 - Keep function sizes small enough for humans to understand them.
 - Unit tests are mandatory for business logic.

## Backend

 - The backend is to be implemented in GoLang.
 - Hexagonal architecture must be used. (ports, adapters, etc.)
 - Node positions and other metadata (data that is relevant for the editor,
   but not for the quests) should be stored in a local sqlite database. But
   this should be implemented with an abstraction layer so that changing
   to another database later on is easier.
 - The backend will either run locally, so no authentication is needed, or
   will be protected by an authenticating proxy server, which will also
   take care of the TLS encryption and certificates.
 - There is no need for user accounts. We do not manage permissions or
   keep track of who edited what.
 - The backend implements the API for the frontend, and also delivers the
   frontend code and assets, so that no additional web servers need to
   be installed.
 - The backend API can provide lists of known conversation partners (which
   are also possible speakers) as well as known items, known factions
   and known resources. Each of these are defined in YAML files, so
   that they can be extended easily.

## Frontend

 - The frontend is a single-page application that runs in the browser.
 - It should be implemeted using the React and React Flow frameworks.
 - But avoid using too many other frameworks, especially if only
   limited functionality is needed from them.
 - It is okay to use Node.js during the build process, but the runtime
   environment should not require Node.js to be installed on the server.
   So, for example, transpiling JSX should be done at build time using
   tools like create-react-app, Webpack or via in-browser solutions like
   Babel Standalone.
 - All assets, such as fonts, graphics, style sheets, javascript files
   etc., must be server by our own backend. No third party links or
   dependencies within the browser. Keep local copies of the required
   files instead.
 - The frontend must validate the quest while it is edited. It should
   show a friendly green checkmark in the top bar if all is well, and
   a list of warnings on the right side if the quest in it's current
   form violates the schema or the "Rules beyond the Schema".
   (from the file "QUEST_NOTES.md")
 - The frontend should make it impossible to create loops. Nonetheless,
   the frontend should check for forbidden loops, as these might
   have been introduced in other ways, bypassing the frontend.

## Look and Feel

The editor should feel modern, uncluttered, clean and fast. When the
user loads the editor, it should show the empty canvas, and a menu
bar on top that has options for creating a new quest file or opening
an existing one - either by selecting it from a drop-down menu or
by typing into a search box, which will limit the available options
while typing for convenience. The top bar also holds a button
to save the quest currently being edited. Saving the quest should
be possible even if there are warnings and the quest is currently
not adhering to all the rules.

The editor should have a tool box on the left side containing the various
node types, so that they can be dragged onto the canvas to create new
node instances.

Each node type comes with different parameters and options. To edit these,
clicking on the node opens an overlay window with the approriate options.

