import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/auth/repositories/user_repository.dart';

final loginControllerProvider =
    AsyncNotifierProvider.autoDispose<LoginController, String?>(
        LoginController.new);

class LoginController extends AutoDisposeAsyncNotifier<String?> {
  @override
  FutureOr<String?> build() {
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
