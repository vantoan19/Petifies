import 'package:grpc/grpc.dart';
import 'package:mobile/src/providers/model_providers.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/providers/storages_provider.dart';
import 'package:retry/retry.dart';
import 'package:riverpod/riverpod.dart';

class RetryUtils {
  static Future<dynamic> RetryRefreshToken(Function fn, Ref ref) async {
    return retry(
      () => fn(),
      maxAttempts: 2,
      delayFactor: Duration.zero,
      retryIf: (e) =>
          e is GrpcError &&
          e.message == "token has expired" &&
          e.code == StatusCode.unauthenticated,
      onRetry: (e) async {
        final secureStorage = ref.read(secureStorageProvider);
        final String? rfToken = await secureStorage.read(key: "refreshToken");
        try {
          final accessToken = (await ref
                  .read(userServiceProvider)
                  .refreshToken(refreshToken: rfToken!))
              .accessToken;
          await secureStorage.write(key: "accessToken", value: accessToken);
        } catch (e) {
          ref.read(myUserProvider.notifier).SetUser(null);
        }
      },
    );
  }
}
