const MEM_DIR = "memories_data";

// Ensure directory exists
if (!VFPrivate.Exists(MEM_DIR)) {
    VFPrivate.CreateDir(MEM_DIR);
}

const getPath = (name) => {
    // Basic sanitization: remove path traversal attempts and ensure .md extension
    const safeName = name.replace(/[\/\\\s]/g, "_");
    return VFPrivate.Join(MEM_DIR, safeName.endsWith(".md") ? safeName : safeName + ".md");
};

const list = () => {
    try {
        const files = VFPrivate.ReadDir(MEM_DIR);
        if (!files || files.length === 0) {
            return "No memories found.";
        }
        // Filter out directories if any, and map to clean names
        const names = files
            .filter(f => !f.IsDir)
            .map(f => f.Name.replace(/\.md$/, ""))
            .sort((a, b) => a.localeCompare(b));

        if (names.length === 0) return "No memories found.";
        return "Memories:\n- " + names.join("\n- ");
    } catch (e) {
        return `Error listing memories: ${e.message}`;
    }
};

const read = (name) => {
    const path = getPath(name);
    if (!VFPrivate.Exists(path)) {
        return `Error: Memory "${name}" does not exist.`;
    }
    try {
        return VFPrivate.ReadStrFile(path);
    } catch (e) {
        return `Error reading memory "${name}": ${e.message}`;
    }
};

const create = (name, content) => {
    const path = getPath(name);
    if (VFPrivate.Exists(path)) {
        return `Error: Memory "${name}" already exists. Use "edit" or "save" to update.`;
    }
    try {
        VFPrivate.WriteStrFile(path, content);
        return `Memory "${name}" created successfully.`;
    } catch (e) {
        return `Error creating memory "${name}": ${e.message}`;
    }
};

const edit = (name, content) => {
    const path = getPath(name);
    if (!VFPrivate.Exists(path)) {
        return `Error: Memory "${name}" does not exist. Use "create" or "save" to create it.`;
    }
    try {
        VFPrivate.WriteStrFile(path, content);
        return `Memory "${name}" updated successfully.`;
    } catch (e) {
        return `Error updating memory "${name}": ${e.message}`;
    }
};

const save = (name, content) => {
    const path = getPath(name);
    try {
        VFPrivate.WriteStrFile(path, content);
        return `Memory "${name}" saved successfully.`;
    } catch (e) {
        return `Error saving memory "${name}": ${e.message}`;
    }
};

const deleteMemory = (name) => {
    const path = getPath(name);
    if (!VFPrivate.Exists(path)) {
        return `Error: Memory "${name}" does not exist.`;
    }
    try {
        VFPrivate.DeleteFile(path);
        return `Memory "${name}" deleted successfully.`;
    } catch (e) {
        return `Error deleting memory "${name}": ${e.message}`;
    }
};

export default {
    list,
    read,
    create,
    edit,
    save,
    delete: deleteMemory
};
