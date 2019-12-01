Feature: ImgCrop service
  As tester of ImgCrop
  I want to check ImgCrop service in work with remote servers

  Scenario: Image in cache
    When Client make get_image request to image_server via ImgCrop service
    Then ImgCrop find image in cache and return them to user

  Scenario: Image server does not exist
    When ImgCrop service make get_image request to image_server1
    Then ImgCrop should return 500 code to User

  Scenario: Image is not found on image server
    When ImgCrop service make get_image request to image_server2
    Then ImgCrop should return 404 code to User

  Scenario: Content type of file is not image
    When ImgCrop service make get_file request to file_server
    Then ImgCrop should return 400 code to User

  Scenario: Remote server return error
    When ImgCrop service make get_image request to image_server and return err_code and error to ImgCrop
    Then ImgCrop should return 503 code to User