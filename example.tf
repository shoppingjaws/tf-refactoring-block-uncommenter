resource "aws_iam_user" "me" {
  name = "my-user"
}
# moved {
  # from = test1
  # to = test2
# }
# removed {
  # from = test3
  # lifecycle {
    # destroy = true
  # }
# }
# import {
  # to = test4
  # id = "foo.bar"
# }