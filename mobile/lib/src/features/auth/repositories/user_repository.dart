import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/models/tokens.dart';
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
  Future<Either<Failure, Map<String, dynamic>>> logIn(
      {required String email, required String password});
  Future<Either<Failure, UserModel>> getMyInfo();
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
        isActivated: false,
      );

      return right(user);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<Either<Failure, Map<String, dynamic>>> logIn({
    required String email,
    required String password,
  }) async {
    try {
      LoginResponse response =
          await _userService.login(email: email, password: password);
      final user = UserModel(
        id: response.user.id,
        email: response.user.email,
        firstName: response.user.firstName,
        lastName: response.user.lastName,
        isActivated: response.user.isActivated,
      );
      final tokens = Tokens(
        sessionId: response.sessionId,
        accessToken: response.accessToken,
        refreshToken: response.refreshToken,
        accessTokenExpiresAt:
            response.accessTokenExpiresAt.seconds.toInt() * 1000,
        refreshTokenExpiresAt:
            response.refreshTokenExpiresAt.seconds.toInt() * 1000,
      );

      Map<String, dynamic> loginResp = {"user": user, "tokens": tokens};

      return right(loginResp);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<Either<Failure, UserModel>> getMyInfo() async {
    try {
      final user = await _userService.getMyInfo();
      UserModel myUser = UserModel(
          id: user.id,
          email: user.email,
          firstName: user.firstName,
          lastName: user.lastName,
          isActivated: user.isActivated,
          followers: 0,
          following: 0,
          countPost: 0,
          bio: "Tao bi gay");
      return right(myUser);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }
}
