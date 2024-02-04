#!/bin/bash

OPENVPN_DIR=/etc/openvpn
LOG_FILE=$OPENVPN_DIR/oath.log
DATE=$(date +"%Y-%m-%d %H:%M:%S")

# VARIABLES
PASS_FILE=$1    # Password file passed by openvpn-server with "auth-user-pass-verify /opt/scripts/auth_client.sh via-file" in server.conf

if [ ! -e "$1" ] || [ ! -s "$1" ]; then
    echo "${DATE} (AUTH SCRIPT) Argument for PASS_FILE either does not exist or is empty." | tee -a $LOG_FILE
    echo "${DATE} (AUTH SCRIPT) There is no auth file. Exit with error." | tee -a $LOG_FILE
    exit 1
else
    echo "${DATE} (AUTH SCRIPT) Argument for PASS_FILE exists and is not empty. It is: $1" | tee -a $LOG_FILE
fi

# Getting user and password passed by external user to OpenVPN server tmp file
user=$(head -1 $PASS_FILE)
pass=$(tail -1 $PASS_FILE)

echo "${DATE} (AUTH SCRIPT) Authentication attempt for user $user" | tee -a $LOG_FILE # echo "$(date) - Password: $pass" | tee -a $LOG_FILE

OATH_DATA_FILE=$OPENVPN_DIR/ccd/${user}
if [ -f "${OATH_DATA_FILE}" ]; then
    echo "${DATE} (AUTH SCRIPT) OATH_DATA_FILE exists and is a regular file: ${OATH_DATA_FILE}" | tee -a $LOG_FILE
else
    echo "${DATE} (AUTH SCRIPT) OATH_DATA_FILE either does not exist or is not a regular file: ${OATH_DATA_FILE}" | tee -a $LOG_FILE
    echo "${DATE} (AUTH SCRIPT) DISABLE AUTH AND ALLOW ACCESS FOR ${user}" | tee -a $LOG_FILE
    exit 0
fi

# Parsing oath.key to getting secret entry, ignore case
key=$(grep -i -m 1 "#2FA_KEY:" ${OATH_DATA_FILE} | cut -d: -f2)
if [ -z "$key" ]; then
    echo "${DATE} (AUTH SCRIPT) OTP KEY IS EMPTY: SKIP OTP CHECKING" | tee -a $LOG_FILE
    code=""
else
    echo "${DATE} (AUTH SCRIPT) OTP KEY IS OK: GENERATE OTP CODE TO VERIFY" | tee -a $LOG_FILE
    # Getting 2FA code with oathtool based on our key, exiting with 0 if match:
    code=$(oathtool --totp=SHA256 ${key})
fi

# Parsing static_pass to getting secret entry, ignore case
static_pass=$(grep -i -m 1 "#STATIC_PASS:" ${OATH_DATA_FILE} | cut -d: -f2)
if [ -z "$static_pass" ]; then
    echo "${DATE} (AUTH SCRIPT) STATIC PASS IS EMPTY: SKIP ITS CHECKING" | tee -a $LOG_FILE
    static_pass=""
else
    echo "${DATE} (AUTH SCRIPT) STATIC PASS IS OK: READY TO VERIFY" | tee -a $LOG_FILE
fi

# Check if we have any code to verify
if [ -z "$key" ] && [ -z "$static_pass" ]; then
  echo "${DATE} (AUTH SCRIPT) OTP KEY and STATIC PASS are empty: nothing to verify" | tee -a $LOG_FILE
  echo "${DATE} (AUTH SCRIPT) DISABLE AUTH AND ALLOW ACCESS FOR ${user}" | tee -a $LOG_FILE
  exit 0
fi

if [ "${static_pass}${code}" = "${pass}" ];
then
    echo "${DATE} (AUTH SCRIPT) Authentication is DONE for user $user" | tee -a $LOG_FILE
    exit 0
else
    echo "FAIL"
fi

# If we make it here, auth hasn't succeeded, don't grant access
echo "${DATE} (AUTH SCRIPT) Authentication failed for user $user" | tee -a $LOG_FILE
exit 1