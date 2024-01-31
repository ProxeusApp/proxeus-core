#!/bin/bash
# In order to reset Proxeus to a pre-defined status, this script
# copies the databases from an origin and copies them to the destination
echo "Restoring database.."
originDataDir="/app/demo/restore_db/*"
destinationDataDir="/data/hosted/"

echo "Removing $destinationDataDir"

rm -R ${destinationDataDir}*

echo "Resetting to initial status from $originDataDir.."
cp -R ${originDataDir} ${destinationDataDir}
echo "Done.."
