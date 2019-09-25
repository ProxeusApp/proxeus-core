#!/bin/bash
if [ `id -u` -ne 0 ]; then
    echo "Execute only as root, Exiting.."
    exit 1
fi

CRON_FILE="/etc/cron.d/demo"
dockerImageToStop="proxeus/proxeus-core:demo"
RESTORE_COMMAND="/usr/bin/docker exec \$(/usr/bin/docker ps -a -q --filter ancestor=$dockerImageToStop --format='{{.ID}}') /app/demo/restore-demo.sh && /usr/bin/docker container restart \$(/usr/bin/docker ps -a -q --filter ancestor=$dockerImageToStop --format='{{.ID}}')"

case "$1" in
    install)
        if [ ! -f $CRON_FILE ]; then
            echo "Cron file for demo not present, creating.."
            touch $CRON_FILE
            chown root:root $CRON_FILE
        fi

        grep -qi "root" $CRON_FILE

        if [ $? != 0 ]; then
            echo "Create cron task to clean environment every 3 days"
            /bin/echo "* * */3 * * root "$RESTORE_COMMAND" >/dev/null 2>&1" >> $CRON_FILE
        fi
        ;;

    clean)
        rm $CRON_FILE
        ;;

    *)
        echo $"Usage: $0 {install|clean}"
        exit 1

esac

