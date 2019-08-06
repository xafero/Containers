#!/bin/sh
find . -name *.sln -exec dotnet restore {} \;
find . -name *.sln -exec dotnet build {} \;
