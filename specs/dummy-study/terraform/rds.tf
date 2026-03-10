resource "aws_db_instance" "main" {
  identifier     = "${var.project}-db"
  engine         = "postgres"
  engine_version = "16"
  instance_class = "db.t4g.micro"

  db_name  = "study"
  username = "postgres"

  allocated_storage     = 20
  max_allocated_storage = 100
  storage_encrypted     = true

  skip_final_snapshot = true

  tags = {
    Project = var.project
  }
}
