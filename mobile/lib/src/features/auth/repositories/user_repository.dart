import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/user_service.dart';

final userRepositoryProvider = Provider(
    (ref) => UserRepository(userService: ref.read(userServiceProvider)));

abstract class IUserRepository {
  Future<Either<Failure, UserModel>> create(
      {required String email,
      required String password,
      required String firstName,
      required String lastName});
  Future<Either<Failure, String>> logIn(
      {required String email, required String password});
}

class UserRepository implements IUserRepository {
  final UserService _userService;

  UserRepository({
    required UserService userService,
  }) : _userService = userService;

  Future<Either<Failure, UserModel>> create(
      {required String email,
      required String password,
      required String firstName,
      required String lastName}) async {
    try {
      User response = await _userService.createUser(
          email: email,
          password: password,
          firstName: firstName,
          lastName: lastName);

      UserModel user = UserModel(
          id: response.id,
          email: response.email,
          firstName: response.firstName,
          lastName: response.lastName,
          isAuthenticated: false,
          isActivated: false);

      return right(user);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<Either<Failure, String>> logIn({
    required String email,
    required String password,
  }) async {
    try {
      LoginResponse response =
          await _userService.login(email: email, password: password);

      return right(response.accessToken);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }
}
