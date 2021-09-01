## **Image Builder**
![Version](https://img.shields.io/badge/version-v0.1.9-blue)
| Artifact | Build Status |
| --- | --- |
| 1-18 OVA | ![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiS0Yvbk1VNW93ZUFxNEFvMC90Ty91d1owRnArN01NczJMalhvbTRYZ09zcWxHNWdCeGdJRFVDU0ZRQXpGQjROVHB6M1o0QzJ4U3V6STU1eVY3TjZYdlhFPSIsIml2UGFyYW1ldGVyU3BlYyI6Ill3VHkxaUV6N2ZaSDlGMTAiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |
| 1-19 OVA | ![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiMjJQV2E2OE5ITVErTEhSTnhGT0Vua081T0lBMHJIalNRM2FPVGhydlVZNERQUHBCVUN3ejU3Ni9Ubng3WjlQYkJrQ3NJSDFwcXczNEdYTWljTW42Uk00PSIsIml2UGFyYW1ldGVyU3BlYyI6ImE4MHNsRGtDd0ttZXlLa0oiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |
| 1-20 OVA | ![BuildStatus](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiazF2R2J2ell0Y0tHT1RnYmF6WXdnRjMwTHMyMTlSVXZnMVoyRytWZ0FDaE5HOU5WejA2VjFzSVNObWlXTjM0eHh2akpBbjgwV0xaTjl5cjFOZlFrZlNNPSIsIml2UGFyYW1ldGVyU3BlYyI6IjhxZTMzVVhZZnR6V0JBOU4iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |
| 1-21 OVA | ![BuildStatus](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoibHJVYmMvSUF0ZlkrVEJsMVBwZU9xLy9ndUZ0U3dGZStpelk2RDRpRTBLQnBrQWNqVkU2TW9qWWI1aFBJM1hpQ1B6TzhaeVduTWdxcE5JeS9XWGhDME5RPSIsIml2UGFyYW1ldGVyU3BlYyI6InVGUS9yandMWmd1cWRsOWciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |

The [Image Builder project](https://github.com/kubernetes-sigs/image-builder) offers a collection of cross-provider Kubernetes virtual machine image building utilities. It can be used to build images intended for use with Kubernetes Cluster API providers. Each provider has its own format of images that it can work with, for example, AMIs for AWS instances, and OVAs for vSphere. The Image Builder project relies on Packer configuration files and Ansible playbooks to build the images and store them in appropriate locations and accounts.

EKS-A CLI project uses these images as the node image when constructing workload clusters for different infrastructure providers like AWS and vSphere.
