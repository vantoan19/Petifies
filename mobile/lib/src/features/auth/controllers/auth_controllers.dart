import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/auth/repositories/user_repository.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'auth_controllers.g.dart';

@Riverpod(keepAlive: false)
class LoginController extends _$LoginController {
  @override
  FutureOr<Map<String, dynamic>?> build() {
    return null;
  }

  Future<void> handle({required String email, required String password}) async {
    final userRepository = ref.read(userRepositoryProvider);

    state = AsyncLoading();
    final result = await userRepository.logIn(email: email, password: password);
    state = result.fold((l) => AsyncValue.error(l, StackTrace.current),
        (r) => AsyncValue.data(r));
  }
}
