import tool from "./tool";

const args = GetArgs();

const main = () => {
    try {
        if (!args.command) {
            return "Error: No command provided. Available commands: list, read, create, edit, save, delete.";
        }

        switch (args.command) {
            case "list":
                return tool.list();
            case "read":
                if (!args.name) return "Error: name is required for read";
                return tool.read(args.name);
            case "create":
                if (!args.name || !args.content) return "Error: name and content are required for create";
                return tool.create(args.name, args.content);
            case "edit":
                if (!args.name || !args.content) return "Error: name and content are required for edit";
                return tool.edit(args.name, args.content);
            case "save":
                if (!args.name || !args.content) return "Error: name and content are required for save";
                return tool.save(args.name, args.content);
            case "delete":
                if (!args.name) return "Error: name is required for delete";
                return tool.delete(args.name);
            default:
                return `Error: Unknown command "${args.command}".`;
        }
    } catch (e) {
        return `Error: ${e.message}`;
    }
};

export default main();
