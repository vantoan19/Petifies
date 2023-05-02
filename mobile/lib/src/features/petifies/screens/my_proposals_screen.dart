import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_controllers.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_proposal_controllers.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/widgets/appbars/my_petifies_appbar.dart';
import 'package:mobile/src/widgets/appbars/only_go_back_appbar.dart';
import 'package:mobile/src/widgets/petifies/petifies.dart';
import 'package:mobile/src/widgets/petifies/proposal.dart';

class MyProposalsScreen extends ConsumerWidget {
  const MyProposalsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);
    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => null,
    );

    if (err != null) {
      return Center(child: Text(err));
    }

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    if (userInfo == null) {
      return Center(child: CircularProgressIndicator());
    }

    final proposals =
        ref.watch(listProposalsByUserIdControllerProvider(userId: userInfo.id));

    return Scaffold(
      appBar: OnlyGoBackAppbar(),
      body: proposals.when(
        data: (data) {
          return Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: Constants.petifiesExpoloreHorizontalPadding,
              vertical: 8,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  "My Proposals",
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Expanded(
                  child: Padding(
                    padding: const EdgeInsets.only(top: 40.0),
                    child: ListView.separated(
                      itemBuilder: (context, index) => ProviderScope(
                        child: ProviderScope(
                          key: ObjectKey(data[index].id),
                          overrides: [
                            proposalProvider.overrideWithValue(data[index])
                          ],
                          child: const PetifiesProposal(),
                        ),
                      ),
                      separatorBuilder: (context, index) => Divider(
                        height: 40,
                        thickness: 0,
                      ),
                      itemCount: data.length,
                    ),
                  ),
                ),
              ],
            ),
          );
        },
        error: (error, stackTrace) => Center(
          child: Text(error.toString()),
        ),
        loading: () => Center(
          child: CircularProgressIndicator(),
        ),
      ),
    );
  }
}
