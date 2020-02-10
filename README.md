# Encrypted file system API

```sh
$ docker run -p 443:443 --name efs-engine efs-engine
> Use `docker exec efs-engine efsctl` to perform commands.

$ alias efsctl="docker exec efs-engine efsctl"

$ efsctl create-fs user-data
> Creating '/efs/user-data'
> Creating files directory '/efs/user-data/files'
> Creating system config directory '/efs/user-data/system'
> Creating AES key '/efs/user-data/system/aes.key'
> Creating filename salt '/efs/user-data/system/filename-salt'

$ efsctl create-user app1 --fs user-data
Copy app1 PGP public key:
<...>
Pub key hash:
    - sha256: <...>
    - sha1: <...>
    - md5: <...>
Use this key for TLS Client Certificate Authentication:
<...>
> Creating user '/efs/user-data/system/users/app1'
> Saving user public key '/efs/user-data/system/users/app1/pub.pgp'
```

```sh
$ docker --name efsc run -v ./keys efs-client

$ docker exec -ti efsc efsc -h serverHost -u priv.pgp -c clientCert.crt
Connected to serverHost...
TLS Client Auth: OK.
PGP User Auth: OK.
Connection to efs server active.

efs> help
post <path>           create file at path with content of stdin
get <path>            get file at path
clear-cache           deletes decrypted local copies of files
efs> post /patient/0/diagnostic.pdf
Paste content:
<...>
Encryption...
Storage...
Done.
efs> clear-cache
`rm -rf /uefs/*` will be performed. Continue ? (y/N) : y
Cache cleared.
efs> stop
shutting down container..
```

```sh
$ efs cp ./Important-Document.pdf drive:/Documents/
$ efs cp --sync --no-delete drive:/Documents /Users/Martin/Documents
$ efs tar --fs drive -o /var/backups/drive.efs.tar.gz.aes -p ./password.txt
```

POST /data/patients/173/prescription.pdf
GET /data/...

Fonctionnement

MS = MicroService

MS1 veut stocker un fichier sensible (ex: prescription uploadée par un utilisateur). Il fait son process dans son fs interne (ex: génération de signature, conversion de type, compression, ...) et ensuite il envoie ses fichiers dans le fs chiffré par API https, avant de recevoir la confirmation de stockage, et de supprimer sa sauvegarde locale. Le contenu reçu par l'API du fs chiffré est d'abord stocké en clair dans le fs interne (non partagé), il est ensuite encodé plusieurs fois puis chiffré en AES256 avant d'être stocké dans un volume persistant. Si la transmission ou le stockage échoue, une erreur est renvoyée à MS1 et MS1 peut re-tenter de stocker le fichier.

Quand MS1 veut récupérer le fichier, il envoie une requête GET sur FS qui authentifie la requête puis regarde s'il dispose d'une version déchiffrée des données dans son fs interne. S'il ne l'a pas, il va la chercher dans le volume partagé et il la déchiffre puis la stocke, ensuite, il renvoie le fichier à MS1.

Authentification: Certificat HTTPs (client cert auth), 2 headers: 1 token basic, 1 auth token.

Quand le FS démarre, il charge la StarterKey depuis sa mémoire qui chiffre les fichiers de clés FILES-AES et JWTKeys.
Le FS écoute sur un port, il fait un auth tls avec un client certificate (hmac) puis une auth par handshake PGP. Il délivre ensuite un BasicAuth valide pour cette instance et un token jwt valable 60 minutes, renouvelable 24 fois.
Les fichiers sont chiffrés en AES avec la clé FILES-AES et sont stockés dans des fichiers nommés par des uuid (référencés dans une bdd boltdb) et arrangés dans un système de fichier optimisé pour s'adapter à une répartition optimale des fichiers selon le fs du volume partagé.
Ex: s'il y a 1 milliard de fichiers à stocker sur un fs unix, les fichier sont répartis dans des sous dossiers de max 65535 fichiers.
Quand il y a des opérations à éxécuter sur le fs, le FS lance une sauvegarde (duplique le dossier racine), puis effectue les modifications sur le système dupliqué pendant qu'il record les modifications sur son fs. Une fois les modifications faites, il commence à écrire sur les 2 fs en même temps puis applique les diff. Avant d'effectuer la transition, il lance un diff général et corrige s'il y a des changements, une fois que c'est bon, il abandonne l'ancien fs et le garde en sauvegarde avant de libérer l'espace.


ex: 65535 > 65535
    3f12/
        de0f
        


