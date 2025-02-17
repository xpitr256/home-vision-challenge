# Checkbox Detection API

This project is a backend solution for detecting checkboxes in images. The API is deployed on Heroku and can be accessed at:

- [Live API Deployment](https://home-vision-challenge-d4dde1803160.herokuapp.com/)
- [Example API Usage](https://home-vision-challenge-d4dde1803160.herokuapp.com//checkbox)

## Version 1.0

- This initial version processes a test image and detects checkboxes.
- Future versions will include a frontend with image upload support.

## Running Locally

To test the API with a different image:

1. Clone this repository:
   ```sh
   git clone https://github.com/xpitr256/home-vision-challenge.git
   cd home-vision-challenge
   ```
2. Replace the test image:
    - Go to the `test/` directory.
    - Replace `test-image.jpg` with another image.
3. Run the application:
   ```sh
   go run main.go
   ```

## API Usage

- Send a request to the `/checkbox` endpoint.
- You can specify an optional `size` parameter (integer between 1 and 200).

### Example Request:

```
GET /checkbox?size=24
```

### Example Response:

```json
{
  "image_name": "test-image.jpg",
  "total_detections": 42,
  "checkboxes": [
    { "x": 156, "y": 96, "status": "unchecked" },
    { "x": 282, "y": 96, "status": "checked" }
  ],
  "size_in_pixels": 24,
  "image_with_checkboxes_url": "/response/image_with_checkboxes.jpg"
   
}
```

## Future Improvements

- Enhance detection accuracy with improved algorithms.
- Support additional image formats.