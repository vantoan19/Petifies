// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_proposal_controllers.dart';
import 'package:mobile/src/models/petifies_proposal.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/appbars/only_go_back_appbar.dart';
import 'package:mobile/src/widgets/petifies/proposal.dart';

class ListProposalsForSessionScreen extends ConsumerWidget {
  final String sessionId;
  const ListProposalsForSessionScreen({
    required this.sessionId,
  });

  Widget renderButtons(PetifiesProposalStatus status) {
    switch (status) {
      case PetifiesProposalStatus.ACCEPTED:
        return ElevatedButton(
          onPressed: () {},
          child: Text("Reject"),
          style: ElevatedButton.styleFrom(
            backgroundColor: Themes.redColor,
            textStyle: TextStyle(color: Themes.blackColor),
            minimumSize: Size(120, 32),
          ),
        );
      case PetifiesProposalStatus.CANCELLED:
        return SizedBox.shrink();
      case PetifiesProposalStatus.REJECTED:
        return SizedBox.shrink();
      case PetifiesProposalStatus.SESSION_CLOSED:
        return SizedBox.shrink();
      case PetifiesProposalStatus.WAITING_FOR_ACCEPTANCE:
        return Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          mainAxisSize: MainAxisSize.max,
          children: [
            ElevatedButton(
                onPressed: () {},
                child: Text(
                  "Accept",
                  style: TextStyle(
                    color: Themes.whiteColor,
                  ),
                ),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Themes.blueColor,
                  minimumSize: Size(120, 32),
                )),
            Padding(
              padding: const EdgeInsets.only(left: 16.0),
              child: ElevatedButton(
                onPressed: () {},
                child: Text("Reject"),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Themes.redColor,
                  minimumSize: Size(120, 32),
                ),
              ),
            ),
          ],
        );
      default:
        return SizedBox.shrink();
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final proposals = ref.watch(
        listProposalsBySessionIdControllerProvider(sessionId: sessionId));

    return Scaffold(
      appBar: const OnlyGoBackAppbar(),
      body: Padding(
        padding: const EdgeInsets.fromLTRB(
          Constants.petifiesExpoloreHorizontalPadding,
          0,
          Constants.petifiesExpoloreHorizontalPadding,
          0,
        ),
        child: proposals.when(
          data: (data) {
            return (data.length > 0)
                ? Expanded(
                    child: ListView.separated(
                      scrollDirection: Axis.horizontal,
                      itemBuilder: (context, index) {
                        return ProviderScope(
                          key: ObjectKey(data[index].id),
                          overrides: [
                            proposalProvider.overrideWithValue(data[index])
                          ],
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const PetifiesProposal(),
                              Padding(
                                padding: const EdgeInsets.only(top: 8.0),
                                child: renderButtons(data[index].status),
                              ),
                            ],
                          ),
                        );
                      },
                      separatorBuilder: (context, index) => Divider(
                        height: 40,
                        thickness: 0,
                      ),
                      itemCount: data.length,
                    ),
                  )
                : const _EmptyItemPlaceholder();
          },
          error: (error, stackTrace) => Center(
            child: Text(error.toString()),
          ),
          loading: () => Center(child: CircularProgressIndicator()),
        ),
      ),
    );
  }
}

class _EmptyItemPlaceholder extends StatelessWidget {
  const _EmptyItemPlaceholder({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(top: 12),
      child: Center(
        child: SizedBox(
          height: 160,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            mainAxisSize: MainAxisSize.min,
            children: [
              SizedBox(
                width: 150,
                height: 150,
                child: Image.asset(Constants.emptyBoxPng),
              ),
              Text("No proposal to show.")
            ],
          ),
        ),
      ),
    );
  }
}
