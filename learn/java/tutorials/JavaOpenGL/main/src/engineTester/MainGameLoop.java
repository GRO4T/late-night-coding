package engineTester;

import entities.*;
import entities.camera.ThirdPersonCamera;
import models.TexturedModel;
import org.lwjgl.opengl.Display;
import org.lwjgl.util.vector.Vector3f;
import renderEngine.*;
import models.RawModel;
import terrains.Terrain;
import textures.ModelTexture;
import textures.TerrainTexture;
import textures.TerrainTexturePack;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;

public class MainGameLoop {
    private static List<Entity> entityList = new ArrayList<>();
    private static List<Terrain> terrainList = new ArrayList<>();

    public static void main(String[] args){
        DisplayManager.createDisplay();
        MasterRenderer renderer = new MasterRenderer();
        Loader loader = new Loader();

        createTerrain(loader);
        generateEntities(loader);

        Light light = new Light(new Vector3f(0, 0, -20), new Vector3f(1, 1, 1));

        RawModel playerRawModel = OBJLoader.loadObjModel("models/person", loader);
        TexturedModel playerTexturedModel = new TexturedModel(playerRawModel,
                new ModelTexture(loader.loadTexture("textures/playerTexture")));
        Player player = new Player(playerTexturedModel, new Vector3f(50, 0, 50), 0, 0, 0, 1);

        ThirdPersonCamera camera = new ThirdPersonCamera(player);

        while(!Display.isCloseRequested()){
            camera.move();
            player.move(getTerrain(player.getPosition().x, player.getPosition().z));

            for (Terrain terrain : terrainList){
                renderer.processTerrain(terrain);
            }
            for (Entity entity : entityList){
                renderer.processEntity(entity);
            }

            renderer.processEntity(player);

            renderer.render(light, camera);
            DisplayManager.updateDisplay();

            System.out.println("frametime: " + DisplayManager.getFrameTimeSeconds());
        }

        renderer.cleanUp();
        loader.cleanUp();
        DisplayManager.closeDisplay();
    }

    private static void generateEntities(Loader loader){
        RawModel treeRawModel = OBJLoader.loadObjModel("models/lowPolyTree", loader);
        RawModel grassRawModel = OBJLoader.loadObjModel("models/grassModel", loader);
        RawModel fernRawModel = OBJLoader.loadObjModel("models/fern", loader);

        ModelTexture treeTexture = new ModelTexture(loader.loadTexture("textures/lowPolyTree"));

        ModelTexture grassTextureAtlas = new ModelTexture(loader.loadTexture("textureAtlases/diffuse"));
        grassTextureAtlas.setHasTransparency(true);
        grassTextureAtlas.setUseFakeLighting(true);
        grassTextureAtlas.setNumberOfRows(4);

        ModelTexture fernTextureAtlas = new ModelTexture(loader.loadTexture("textureAtlases/fern"));
        fernTextureAtlas.setHasTransparency(true);
        fernTextureAtlas.setUseFakeLighting(true);
        fernTextureAtlas.setNumberOfRows(2);

        TexturedModel treeStaticModel = new TexturedModel(treeRawModel, treeTexture);
        TexturedModel grassStaticModel = new TexturedModel(grassRawModel, grassTextureAtlas);
        TexturedModel fernStaticModel = new TexturedModel(fernRawModel, fernTextureAtlas);

        Random generator = new Random();
        for (int i = 0; i < 200; i++){
            createEntityAtRandomPos(treeStaticModel, generator, 0);
            for (int j = 0; j < 10; j++){
                createEntityAtRandomPos(grassStaticModel, generator, generator.nextInt(9));
            }
            for (int j = 0; j < 5; j++){
                createEntityAtRandomPos(fernStaticModel, generator, generator.nextInt(4));
            }
        }
    }

    private static void createTerrain(Loader loader){
        TerrainTexture backgroundTexture = new TerrainTexture(loader.loadTexture("textures/grassy2"));
        TerrainTexture rTexture = new TerrainTexture(loader.loadTexture("textures/mud"));
        TerrainTexture gTexture = new TerrainTexture(loader.loadTexture("textures/grassFlowers"));
        TerrainTexture bTexture = new TerrainTexture(loader.loadTexture("textures/path"));

        TerrainTexturePack texturePack = new TerrainTexturePack(backgroundTexture, rTexture,
                gTexture, bTexture);
        TerrainTexture blendMap = new TerrainTexture(loader.loadTexture("blendMaps/blendMap"));

        terrainList.add(new Terrain(0, 0, loader, texturePack, blendMap, "heightMaps/heightMap"));
        terrainList.add(new Terrain(1, 0, loader, texturePack, blendMap, "heightMaps/heightMap"));
        terrainList.add(new Terrain(1, 1, loader, texturePack, blendMap, "heightMaps/heightMap"));
        terrainList.add(new Terrain(0, 1, loader, texturePack, blendMap, "heightMaps/heightMap"));
    }

    private static Terrain getTerrain(float worldX, float worldZ){
        for (Terrain terrain : terrainList){
            if (worldX >= terrain.getX() && worldX < terrain.getX() + terrain.getSIZE() &&
                    worldZ >= terrain.getZ() && worldZ < terrain.getZ() + terrain.getSIZE()){
                return terrain;
            }
        }
        return null;
    }

    private static void createEntityAtRandomPos(TexturedModel model, Random generator, int textureID){
        float boundX = 1600;
        float boundZ = 1600;
        float x = generator.nextFloat() * boundX;
        float z = generator.nextFloat() * boundZ;
        float y = getTerrain(x, z).getHeightOfTerrain(x, z);

        entityList.add(new Entity(model, textureID, new Vector3f(x, y, z), 0, 0, 0, 1));
    }
}

