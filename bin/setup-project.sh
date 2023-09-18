#!/bin/ash

# If the argument is empty, then exit with error message
if [ -z "$1" ]; then
	echo "Client directory name is required"
	exit 1
fi
CLIENT_DIR_NAME="$1"
if [ -z "$2" ]; then
	echo "Project directory name is required"
	exit 1
fi
PROJECT_DIR_NAME="$2"
if [ -z "$3" ]; then
	echo "Runtime is required"
	exit 1
fi
RUNTIME="$3"
if [ -z "$4" ]; then
	echo "Database is required"
	exit 1
fi
DATABASE="$4"
if [ -z "$5" ]; then
	echo "Owner email is required"
	exit 1
fi
OWNER_EMAIL="$5"

DOCUMENT_ROOT="/usr/local/apache2/htdocs/"

TARGET_DIR="${DOCUMENT_ROOT}${CLIENT_DIR_NAME}/${PROJECT_DIR_NAME}"

# If the target directory already exists, then exit with error message
if [ -d "${TARGET_DIR}" ]; then
	echo "Target directory already exists. '${TARGET_DIR}'"
	exit 1
fi
# create the target directory
mkdir -p "${TARGET_DIR}"

# copy the template files to the target directory
cp "${DOCUMENT_ROOT}template/index.html" "${TARGET_DIR}/index.html"

# replace the placeholders with the actual values
sed -i "s/%CLIENT%/${CLIENT_DIR_NAME}/g" "${TARGET_DIR}/index.html"
sed -i "s/%PROJECT%/${PROJECT_DIR_NAME}/g" "${TARGET_DIR}/index.html"
sed -i "s/%OWNER%/${OWNER_EMAIL}/g" "${TARGET_DIR}/index.html"
sed -i "s/%DATABASE%/${DATABASE}/g" "${TARGET_DIR}/index.html"
sed -i "s/%RUNTIME%/${RUNTIME}/g" "${TARGET_DIR}/index.html"

echo "The project has been created."
exit
