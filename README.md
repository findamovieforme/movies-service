# movies-service

## Local development

### Python (recommendation predictor)

The `/movies/recommendations` endpoint calls a local Python script that needs `joblib` and `pandas`. Install them once:

```bash
pip install -r helpers/requirements.txt
```

Or with a venv (recommended):

```bash
python3 -m venv .venv
source .venv/bin/activate   # or .venv\Scripts\activate on Windows
pip install -r helpers/requirements.txt
```

Run the service from the **project root** so `helpers/predictor.py` and `recommendation-model/` are found.

### Docker / AWS

The Docker image includes Python 3, the predictor dependencies, and the recommendation model. No extra setup is needed in AWS when using this image.

