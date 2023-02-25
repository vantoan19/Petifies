import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:grpc/service_api.dart';
import 'package:mobile/src/providers/storages_provider.dart';

class AuthInterceptor extends ClientInterceptor {
  final Ref _ref;

  AuthInterceptor({required ref}) : _ref = ref;

  FutureOr<void> _provider(Map<String, String> metadata, String uri) async {
    final secureStorage = _ref.read(secureStorageProvider);
    var accessToken = await secureStorage.read(key: 'accessToken');
    metadata['authorization'] = 'Bearer $accessToken';
  }

  @override
  ResponseStream<R> interceptStreaming<Q, R>(
      ClientMethod<Q, R> method,
      Stream<Q> requests,
      CallOptions options,
      ClientStreamingInvoker<Q, R> invoker) {
    return super.interceptStreaming(method, requests,
        options.mergedWith(CallOptions(providers: [_provider])), invoker);
  }

  @override
  ResponseFuture<R> interceptUnary<Q, R>(ClientMethod<Q, R> method, Q request,
      CallOptions options, ClientUnaryInvoker<Q, R> invoker) {
    return super.interceptUnary(method, request,
        options.mergedWith(CallOptions(providers: [_provider])), invoker);
  }
}
