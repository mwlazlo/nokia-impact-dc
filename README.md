## Step 1: build lwm2mclient
```bash
git clone https://github.com/eclipse/wakaama
cd wakaama/examples/shared
git clone https://git.eclipse.org/r/tinydtls/org.eclipse.tinydtls tinydtls
cd ../client
perl -pi -e 's/"Enable DTLS" OFF/"Enable DTLS" ON/' CMakeLists.txt
cmake .
make
file lwm2mclient
echo output should be: lwm2mclient: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, for GNU/Linux 3.2.0, BuildID[sha1]=879a4fbbce0c24e7b2aa94fad1ac7676dc0a3995, not stripped
```

## Step 2: Setup google cloud 
- Create a new project: https://console.cloud.google.com
- Install SDK: https://cloud.google.com/sdk/install
- Download service account private key: https://console.firebase.google.com/u/0/project/{{YOUR-PROJECT-ID}}/settings/serviceaccounts/adminsdk
- Setup the SDK:
`gcloud init`

## Step 3: Deploy backend
```bash
git clone https://github.com/mwlazlo/nokia-impact-dc
cd nokia-impact-dc/
cd cmd/nokia-impact-dc
cp ~/Downloads/service-account-xyz123.json .
cat > config.json <<EOF
{ 
  "CallbackUsername": "nokia",
  "CallbackPassword": "Nokia@20190218",
  "ImpactUsername": "Ubiik.TH",
  "ImpactPassword": "Ubiik@19",
  "ImpactBaseURL": "https://impact.idc.nokia.com",
  "ImpactGroup": "APJ.JAPAN.Rakuten",
  "GoogleAuthFile": "ubiik-auth.json",
  "ListenPort": "8080"
}
EOF
gcloud app deploy
```

## Step 4: Test lwm2mclient

- Watch logfiles: gcloud logs tail
- Spin up a client: `./lwm2mclient -h impact.idc.nokia.com -n $(uuid) -p 30001 -4 -i mypsid -s 123456789abcdef -c`
- Watch the data come in: https://console.firebase.google.com 


