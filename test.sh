#!/bin/bash

ls release/* | xargs -I FILE basename FILE
