# auditcleanup

This repository provides a simplified and more secure version of the audit cleanup container for the Couchbase Autonomous Operator.
No guarantees are made for security accreditation but best practices have been followed:

* Scratch based image (not for Red Hat Container Catalog due to their requirements so using UBI there).
* Programmatic Golang executable that only removes rotated audit logs.
* Security and image scanning carried out during builds to check for CVEs and other issues.