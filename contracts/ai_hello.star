# AIVM Starlark Smart Contract Example
# This Python-like contract orchestrates AI model calls on-chain.
#
# Built-in functions available in every contract:
#   request_ai(model_id, input_data)  → returns execution_id (string)
#   get_execution(execution_id)       → returns dict with keys:
#                                         status, output_hash, proof, model_id

def main(args):
    """
    Main entry point. Called when the contract is executed on-chain.
    args is a dict of string key/value pairs passed by the caller.
    """
    prompt = args.get("prompt", "default prompt")
    model  = args.get("model",  "default-classifier")

    # Step 1: Trigger an AI inference via the AIVM dual-path router
    exec_id = request_ai(model, prompt)

    # Step 2: Optionally inspect the execution record
    # (status will be EXECUTED_ON_CHAIN or PENDING_OFF_CHAIN)
    record = get_execution(exec_id)

    if record == None:
        return "Error: execution record not found for " + exec_id

    status = record["status"]

    # Step 3: Return a meaningful result string
    return "AI execution queued! ID=" + exec_id + " status=" + status
