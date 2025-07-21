import sys
import json
import os
import traceback

import joblib
import pandas as pd


def debug_env(model_dir: str) -> None:
    """
    Extra diagnostics printed to stderr to help understand failures
    in the container when loading the model.
    """
    try:
        import sklearn  # type: ignore[import]

        print(
            json.dumps(
                {
                    "debug": "sklearn_available",
                    "sklearn_version": getattr(sklearn, "__version__", "unknown"),
                }
            ),
            file=sys.stderr,
        )
    except Exception as e:  # noqa: BLE001
        print(
            json.dumps(
                {
                    "debug": "sklearn_import_failed",
                    "error": str(e),
                }
            ),
            file=sys.stderr,
        )

    # Basic filesystem checks
    files = {}
    for name in ("tfidf_vectorizer.joblib", "knn_model.joblib", "movies.csv"):
        path = os.path.join(model_dir, name)
        files[name] = {"exists": os.path.exists(path), "path": path}

    print(
        json.dumps(
            {
                "debug": "model_dir_info",
                "model_dir": model_dir,
                "pwd": os.getcwd(),
                "files": files,
            }
        ),
        file=sys.stderr,
    )


def main() -> None:
    # Load model artifacts (path relative to this script so it works from any CWD)
    _script_dir = os.path.dirname(os.path.abspath(__file__))
    model_dir = os.path.join(_script_dir, "..", "recommendation-model")

    # Extra debugging about environment and model files
    debug_env(model_dir)

    try:
        vectorizer = joblib.load(f"{model_dir}/tfidf_vectorizer.joblib")
        knn_model = joblib.load(f"{model_dir}/knn_model.joblib")
        movies = pd.read_csv(f"{model_dir}/movies.csv")
    except Exception as e:  # noqa: BLE001
        # Return a structured error so the Go service can surface it cleanly
        print(
            json.dumps(
                {
                    "error": f"failed to load model artifacts: {e}",
                    "traceback": traceback.format_exc(),
                }
            )
        )
        sys.exit(1)

    # Read input JSON from stdin
    try:
        input_data = json.load(sys.stdin)
    except Exception as e:  # noqa: BLE001
        print(json.dumps({"error": f"invalid JSON input: {e}"}))
        sys.exit(1)

    title = input_data.get("title")
    n = int(input_data.get("n", 5))

    if title not in movies["title"].values:
        print(json.dumps({"error": f"Movie '{title}' not found in the dataset."}))
        sys.exit(0)

    idx = movies[movies["title"] == title].index[0]
    movie_features = movies.loc[idx, "combined_features"]
    movie_vec = vectorizer.transform([movie_features])
    distances, indices = knn_model.kneighbors(movie_vec, n_neighbors=n + 1)

    recommended_indices = indices.flatten()[1:]
    recommended_titles = movies["title"].iloc[recommended_indices].tolist()

    print(json.dumps({"recommendations": recommended_titles}))


if __name__ == "__main__":
    main()
