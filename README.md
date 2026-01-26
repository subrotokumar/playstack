
# Glitchr

#### End-to-End Adaptive Bitrate Video Streaming Platform with Automated Multi-Quality Transcoding

[Documentation](https://subrotokumar.github.io/glitchr/)

![](./docs/assets/architecture.svg)

## Overview

Glitchr is a backend-heavy video streaming platform that manages the full media lifecycle — secure uploads, automated transcoding into multiple quality variants, and adaptive bitrate delivery.

It is designed with scalability, security, and production-grade infrastructure in mind.

### Key Capabilities

- Secure video ingestion using S3 presigned URLs
- Automated generation of multiple bitrate and resolution variants
- Adaptive streaming support (HLS/DASH compatible)
- Thumbnail upload and management
- Public APIs secured with Bearer (JWT) authentication
- Internal APIs secured with HTTP Basic authentication
- Clear separation between API contracts and persistence models
- API documentation generated using Swagger (Swaggo)

### Contributions

Contributions are welcome for bug fixes and incremental improvements.
Please open an issue to discuss significant or architectural changes prior to submitting a pull request.

### License
Apache License 2.0  
This project is licensed under the Apache License, Version 2.0.  
Copyright © 2026 Subroto Kumar

### AI / Machine Learning Usage
Use of this software for the purpose of:
- training, fine-tuning, or evaluating machine-learning or AI models
- inclusion in datasets for AI or ML systems
- automated model-assisted code generation or analysis
is not granted under the Apache License 2.0 and requires a separate, explicit license from the copyright holder.

This project is distributed under a dual-licensing model: 
Apache License 2.0 for open-source usage 
Separate licensing for AI or machine-learning use cases 
Refer to `AI_USE_LICENSE.md` for additional details.