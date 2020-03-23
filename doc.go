/*

https://drive.google.com/open?id=1HVJTSiav-OfLfm7yqHgokI67B2_YhJ8N

cerbaes is an example package for Cerberus project
Cerberus handles accounts registrations - holding multiple account data fields,
and allows saving documents data and scanned copies of the documents for each account
each document handles also different versions

peer database - CouchDB
documents location - IPFS network

Accounts:

account creation:
application level code inside app/ receives data from the frontend
unique (private) account id is created using bson.NewObjectId() and is saved in hex
public ID is an Merkle-Damgard MD5 hash of the private ID

account data bytes are encrypted with AES256-CGM algorithm, using a randomly created 32-bit key
and saved in peer database - CounchDB
newly created key is returned

account data update:
when account data fields are updated - data is extracted from peer database,
decrypted using the provided key, updated, encrypted again and send to the database
updating the record under the same publicID
note: account data updates are currently implemented inside the chaincode for faster performance

Documents:

document creation:
Document creations is done per specific account, its data is extracted from peer DB providing the key,
decrypted, document data is added inside the account record,
document content is added to IPFS network and its IPFS hash location is added to the
account document details as well
then it is encrypted and sent to peer database, updating the records under the same PublicID
this logic is implemented on application level, because of keeping chaincode code
as simple as possible

Document data inside account structure is encrypted again with the account
ecnryption/decryption key - AES256-GCM

each document version is saved under specific hash corresponding to the account, document name and version
content is encrypted AES256-CGM, using different key
the 256-bit cipherkey is signed with public key from RSA-1024 pair generated for each document version
both RSA pair and the cipherkey are saved inside temporary directory
this is a temporary solution inside the project - they are supposed to be handled by the frontend


*/
package cerbaes
