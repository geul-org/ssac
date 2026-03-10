@scenario
Feature: Create and cancel a reservation

  Scenario: Full reservation lifecycle
    Given POST Login {"Email": "user@test.com", "Password": "Pass1234!"} → token
    When POST CreateReservation {"RoomID": 1, "StartAt": "2025-06-01T00:00:00Z", "EndAt": "2025-06-03T00:00:00Z"} → reservation
    Then status == 200
    And response.reservation exists
    And GET GetReservation {"ReservationID": reservation.ID} → detail
    And response.reservation exists
    When PUT CancelReservation {"ReservationID": reservation.ID}
    Then status == 200
