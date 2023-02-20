class User {
  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final bool isAuthenticated;
  final bool isActivated;

  User(this.id, this.email, this.firstName, this.lastName, this.isAuthenticated,
      this.isActivated);
}
