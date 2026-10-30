# Authoritative profile-picture storage

OAuth forwards uploads to f-image and reports success only after HTTP200 and
`{"status":"success","worked":true}`. Failed transport/storage returns503;
a malformed or negative success acknowledgement returns502. There is no local
OAuth write or read fallback. Reads use f-image with a30-second timeout and the
incoming request context. Existing UI error handling keeps the edit modal and
selected file open for a deliberate retry. A failed/timeout acknowledgement can
mean storage completed but confirmation was lost; no automatic write retry occurs.

f-image remains one replica on Pi5 backed by its retained PVC. This change makes
OAuth replicas stateless for pictures; it does not make f-image highly available.
Provider-imported pictures already use f-image rather than local OAuth files.

## Deployment prerequisites and existing pictures

This source change is not deployed. Running Minikube must remain untouched until
owner authorizes its rollout/migration. Do not replace existing OAuth pods first:
legacy local-only pictures might exist only in their writable container layers.

During the authorized maintenance phase:

1. Fence uploads and inventory every current OAuth pod's `public/images` files
   plus the entire f-image dynamic-data PVC, including profiles and galleries.
   Capture private copies with source pod UID/path, size, timestamp and SHA256;
   retain originals. Do not assume that an image exists centrally because GET
   returns200: f-image can generate a default avatar on reads.
2. Make an encrypted, independently retained archive of both sets and verify a
   restore into an isolated directory. Review filename/user-ID mappings and
   conflicts. A local original and centrally resized PNG have different hashes;
   hash inequality does not identify which version is newer. Resolve conflicts
   explicitly and never blindly overwrite central pictures.
3. Import only approved missing/selected legacy pictures through f-image's upload
   API while other writers remain fenced; verify acknowledgements and resulting
   valid image content. Keep a per-user import journal and preserved originals.
4. Capture and restore-check the final authoritative f-image PVC backup. The
   migration preparation now includes the `f-image-data` native snapshot hook;
   there is no shared OAuth upload volume or directory backup hook to provision.
5. Roll out OAuth, verify upload/read across both replicas, simulate image-service
   failure and check an error with no success/local file, then enable uploads.

Independent backup destination, legacy import and destination acceptance remain
open. Existing dated SQL/SQLite restoration does not prove picture recovery.
