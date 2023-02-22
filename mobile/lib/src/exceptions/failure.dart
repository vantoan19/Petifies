class Failure implements Exception {
  final String? message;

  const Failure([this.message]);

  @override
  String toString() {
    return 'Failure: $message';
  }
}
