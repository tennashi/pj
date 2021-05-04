# Project management specification
Tags: project

## The user can create a project without specifying a project name

The user can create a project without specifying a project name to create a project in the current directory name.

* Move to the directory that the user wants to specify as the project name
* Create a project without specifying a project name
* Show the details of the project

## The user can create a project with specifying the project name

The user can create a project with specifying the project name as "awesome".

* Create a project with specifying the project name as "awesome"
* Show the details of the project "awesome"

## If the user creates a project with a name that already exists, nothing will happen

* Create a project with specifying the project name as "awesome"
* Show the details of the project "awesome"
* Make sure there is only "awesome" in the list of projects
* Create a project with specifying the project name as "awesome"
* Show the details of the project "awesome"
* Make sure there is only "awesome" in the list of projects

## The user can change the current project

* Create a project with specifying the project name as "awesome-1"
* Create a project with specifying the project name as "awesome-2"
* Check that the current project is "awesome-2"
* Change the current project to "awesome-1"
* Check that the current project is "awesome-1"

## The user can add the workspace to the current project

* Create a project without specifying a project name
* Create a new directory with the name "new-workspace"
* Move to the directory "new-workspace"
* Add the current directory as a workspace to the current project
* Check that the current project has "new-workspace" as a workspace

## If the user add a workspace with the path that already exists, nothing will happen

* Create a project without specifying a project name
* Add the current directory as a workspace to the current project
* Check that the current project has only one workspace
