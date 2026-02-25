import os

from flask import Flask, jsonify, request
import joblib
import pandas as pd


app = Flask(__name__)


_SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
_MODEL_DIR = os.path.join(_SCRIPT_DIR, "..", "recommendation-model")

# Load model artifacts once at startup
_VECTORIZER = joblib.load(os.path.join(_MODEL_DIR, "tfidf_vectorizer.joblib"))
_KNN_MODEL = joblib.load(os.path.join(_MODEL_DIR, "knn_model.joblib"))
_MOVIES = pd.read_csv(os.path.join(_MODEL_DIR, "movies.csv"))


@app.get("/health")
def health() -> tuple[dict, int]:
  return {"status": "ok"}, 200


@app.post("/predict")
def predict() -> tuple[dict, int]:
  data = request.get_json(force=True) or {}
  title = data.get("title")
  n = int(data.get("n", 5))

  if not title:
    return {"error": "missing 'title'"}, 400

  if title not in _MOVIES["title"].values:
    return {"error": f"Movie '{title}' not found in the dataset."}, 200

  idx = _MOVIES[_MOVIES["title"] == title].index[0]
  movie_features = _MOVIES.loc[idx, "combined_features"]
  movie_vec = _VECTORIZER.transform([movie_features])
  _, indices = _KNN_MODEL.kneighbors(movie_vec, n_neighbors=n + 1)

  recommended_indices = indices.flatten()[1:]
  recommended_titles = _MOVIES["title"].iloc[recommended_indices].tolist()

  return {"recommendations": recommended_titles}, 200


if __name__ == "__main__":
  app.run(host="0.0.0.0", port=5000)

